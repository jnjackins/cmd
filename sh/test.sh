#!./sh
# run go build before running this test

echo beginning test
echo

echo command-line arguments: '['^$*^']'

echo '''   quoted    string    with    whitespace    '''

echo setting foo '=' bar
foo = bar
echo the value of '$'^foo is $foo

echo writing to a file called test.^$pid
file = test.^$pid
echo reading back from the file >$file
cat <$file
rm $file

echo
echo done testing

